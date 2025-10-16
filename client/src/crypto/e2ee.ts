// Minimal E2EE key generation using Web Crypto.
// Identity: ECDSA P-256 (for signing)
// Signed PreKey: ECDH P-256 public key, signed by identity private key
// One-time PreKeys: ECDH P-256 public keys

import { e2eeAPI } from '../api/client';

function toBase64Url(data: ArrayBuffer): string {
  const b64 = btoa(String.fromCharCode(...new Uint8Array(data)));
  return b64.replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/g, '');
}

async function exportPublicKey(key: CryptoKey, format: 'raw' | 'spki' = 'spki'): Promise<string> {
  const buf = await crypto.subtle.exportKey(format, key);
  return toBase64Url(buf);
}

async function exportPrivateKey(key: CryptoKey): Promise<string> {
  const buf = await crypto.subtle.exportKey('pkcs8', key);
  return toBase64Url(buf);
}

async function signBytes(privateKey: CryptoKey, data: ArrayBuffer): Promise<string> {
  const sig = await crypto.subtle.sign({ name: 'ECDSA', hash: 'SHA-256' }, privateKey, data);
  return toBase64Url(sig);
}

function getOrCreateDeviceId(): string {
  const existing = localStorage.getItem('device_id');
  if (existing) return existing;
  const dev = crypto.randomUUID();
  localStorage.setItem('device_id', dev);
  return dev;
}

export async function ensurePublishedKeyBundle(): Promise<void> {
  const deviceId = getOrCreateDeviceId();
  const publishedMarker = localStorage.getItem(`e2ee_published_${deviceId}`);
  if (publishedMarker === '1') return; // already published for this device

  // Generate identity (ECDSA) if not present
  let identityPrivB64 = localStorage.getItem('e2ee_identity_priv');
  let identityPubB64 = localStorage.getItem('e2ee_identity_pub');

  let identityKeys: CryptoKeyPair | null = null;
  if (!identityPrivB64 || !identityPubB64) {
    identityKeys = await crypto.subtle.generateKey(
      { name: 'ECDSA', namedCurve: 'P-256' },
      true,
      ['sign', 'verify']
    );
    identityPrivB64 = await exportPrivateKey(identityKeys.privateKey);
    identityPubB64 = await exportPublicKey(identityKeys.publicKey, 'spki');
    localStorage.setItem('e2ee_identity_priv', identityPrivB64);
    localStorage.setItem('e2ee_identity_pub', identityPubB64);
  }

  // Import identity private for signing
  const identityPrivate = await crypto.subtle.importKey(
    'pkcs8',
    Uint8Array.from(atob(identityPrivB64!), (c) => c.charCodeAt(0)).buffer,
    { name: 'ECDSA', namedCurve: 'P-256' },
    false,
    ['sign']
  );

  // Generate Signed PreKey (ECDH key, sign its public part)
  const spk = await crypto.subtle.generateKey(
    { name: 'ECDH', namedCurve: 'P-256' },
    true,
    ['deriveBits']
  );
  const spkPubSpki = await crypto.subtle.exportKey('spki', spk.publicKey);
  const spkSig = await signBytes(identityPrivate, spkPubSpki);
  const spkPubB64 = toBase64Url(spkPubSpki);
  // Store signed prekey private for decryption later
  const spkPrivPkcs8 = await crypto.subtle.exportKey('pkcs8', spk.privateKey);
  localStorage.setItem('e2ee_spk_priv', toBase64Url(spkPrivPkcs8));

  // Generate a small pool of One-Time PreKeys (ECDH pub only)
  const oneTime: { key_id: number; public_key: string }[] = [];
  const count = 20;
  for (let i = 1; i <= count; i++) {
    const otk = await crypto.subtle.generateKey(
      { name: 'ECDH', namedCurve: 'P-256' },
      true,
      ['deriveBits']
    );
    const otkPub = await exportPublicKey(otk.publicKey, 'spki');
    oneTime.push({ key_id: i, public_key: otkPub });
  }

  // Publish bundle (only public materials)
  await e2eeAPI.publishDeviceKeys({
    device_id: deviceId,
    identity_key_public: identityPubB64!,
    signed_prekey_public: spkPubB64,
    signed_prekey_signature: spkSig,
    one_time_prekeys: oneTime,
  });

  // Mark as published for this device
  localStorage.setItem(`e2ee_published_${deviceId}`, '1');
}

async function importSpkiToCryptoKey(spkiB64Url: string, algo: 'ECDH' | 'ECDSA', keyUsages: KeyUsage[]): Promise<CryptoKey> {
  const b = atob(spkiB64Url.replace(/-/g, '+').replace(/_/g, '/'));
  const buf = new Uint8Array([...b].map(c => c.charCodeAt(0))).buffer;
  return crypto.subtle.importKey('spki', buf, { name: algo, namedCurve: 'P-256' }, true, keyUsages);
}

async function importPkcs8ToCryptoKey(pkcs8B64Url: string, algo: 'ECDH' | 'ECDSA', keyUsages: KeyUsage[]): Promise<CryptoKey> {
  const b = atob(pkcs8B64Url.replace(/-/g, '+').replace(/_/g, '/'));
  const buf = new Uint8Array([...b].map(c => c.charCodeAt(0))).buffer;
  return crypto.subtle.importKey('pkcs8', buf, { name: algo, namedCurve: 'P-256' }, false, keyUsages);
}

export async function encryptForBundle(bundle: {
  identity_key_public: string;
  signed_prekey_public: string;
  one_time_prekey?: { key_id: number; public_key: string } | null;
}, plaintext: string): Promise<{ ciphertext: string; nonce: string; alg: string; ephemeral_pub: string; }>
{
  // Ephemeral ECDH key pair for the sender
  const eph = await crypto.subtle.generateKey({ name: 'ECDH', namedCurve: 'P-256' }, true, ['deriveBits']);
  const recipientECDHPublic = await importSpkiToCryptoKey(bundle.one_time_prekey?.public_key || bundle.signed_prekey_public, 'ECDH', []);
  const sharedBits = await crypto.subtle.deriveBits({ name: 'ECDH', public: recipientECDHPublic }, eph.privateKey, 256);
  const aesKey = await crypto.subtle.importKey('raw', sharedBits, { name: 'AES-GCM', length: 256 }, false, ['encrypt']);
  const nonce = crypto.getRandomValues(new Uint8Array(12));
  const encoded = new TextEncoder().encode(plaintext);
  const ctBuf = await crypto.subtle.encrypt({ name: 'AES-GCM', iv: nonce }, aesKey, encoded);
  const ephSpki = await crypto.subtle.exportKey('spki', eph.publicKey);
  return {
    ciphertext: toBase64Url(ctBuf),
    nonce: toBase64Url(nonce.buffer),
    alg: 'ECDH-P256+AES-GCM',
    ephemeral_pub: toBase64Url(ephSpki),
  };
}

export async function tryDecryptMessage(message: { ciphertext?: string; nonce?: string; alg?: string; ephemeral_pub?: string; content?: string; }): Promise<string | null> {
  if (!message.ciphertext || !message.nonce || !message.ephemeral_pub) {
    return message.content ?? null;
  }
  const spkPrivB64 = localStorage.getItem('e2ee_spk_priv');
  if (!spkPrivB64) return null;
  try {
    const recipientECDHPrivate = await importPkcs8ToCryptoKey(spkPrivB64, 'ECDH', ['deriveBits']);
    const ephPub = await importSpkiToCryptoKey(message.ephemeral_pub, 'ECDH', []);
    const sharedBits = await crypto.subtle.deriveBits({ name: 'ECDH', public: ephPub }, recipientECDHPrivate, 256);
    const aesKey = await crypto.subtle.importKey('raw', sharedBits, { name: 'AES-GCM', length: 256 }, false, ['decrypt']);
    const nonceBytes = Uint8Array.from(atob(message.nonce.replace(/-/g, '+').replace(/_/g, '/')), c => c.charCodeAt(0));
    const ctBytes = Uint8Array.from(atob(message.ciphertext.replace(/-/g, '+').replace(/_/g, '/')), c => c.charCodeAt(0));
    const ptBuf = await crypto.subtle.decrypt({ name: 'AES-GCM', iv: nonceBytes }, aesKey, ctBytes);
    return new TextDecoder().decode(ptBuf);
  } catch {
    return null;
  }
}


