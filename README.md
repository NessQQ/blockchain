# Transaction-Aware Proof of Work (PoW) Blockchain Engine

An institutional-quality, educational Go blockchain implementation designed to demonstrate core cryptographic consensus principles. This platform integrates asymmetric key pairs, cryptographic wallet infrastructure, digital transaction signing/verification gates, and a deterministic SHA-256 Proof of Work mining engine.

---

## 🏛️ Comprehensive Structural Data Flow Layout

Data in this system moves across strict, decoupled architectural boundaries to guarantee user identity security, transaction authenticity, and immutable ledger integrity.

```unset
 [PHASE 1: IDENTITY REFS]
  wallet.go ───> Generates Asymmetric Private/Public Keys
                      │
                      ▼
 [PHASE 2: TX SIGNING LOCK]
  wallet.go ───> Packs fields into data strings ───> Signs via Private Key ───> Injects "Signature" string into Transaction
                                                                                         │
                                                                                         ▼
 [PHASE 3: CONSENSUS SECURITY GATE]                                                      │
  blockchain.go ───> Intercepts data blocks ───> Loops over user items ───────────────────┤
                          │                                                              │
                          ├───> Skips Coinbase elements (System tokens require no keys)   │
                          │                                                              │
                          └───> VerifyTransaction(tx, pubKey, tx.Signature) <────────────┘
                                     │
                                     ▼ If signatures match securely:
 [PHASE 4: IMMUTABLE LEDGER REFINERY]
  block.go ───> Serializes complete data strings (Fields + Signatures together)
                     │
                     ▼
  proof.go ───> Brute-forces Nonce loop values until matching hash limits are broken ───> Saved into BlockChain arrays
```

---

## ⚡ Core Technical & Cryptographic Pillars

### 1. Asymmetric Wallet Infrastructure & Identity
A wallet is a secure cryptographic container that manages user identities on the network utilizing **RSA public/private key pairs**.
*   **The Private Key:** A strictly guarded secret used to generate unique digital signatures that cryptographically lock and authorize outgoing transfers.
*   **The Public Key:** Shared freely within the network to serve as a user's address and allow other participating nodes to verify the authenticity of their digital signatures.
*   **Transaction Signing:** Outbound transactions are serialized into a strict string pattern (`"Sender->Receiver:Amount(Coinbase)"`) and signed using the private key. This appends a unique base64 `Signature` string directly to the transaction payload.

### 2. Consensus Security Gate (Verification)
Before any transaction is allowed to enter a block or contaminate the ledger array history, it must clear the consensus gate inside `blockchain.go`.
*   The system uses the sender's public key to decrypt the transaction signature string.
*   If an attacker attempts a "Man-in-the-Middle" attack by modifying the receiver or the amount after the user signed it, the calculated cryptographic hash will no longer match the signature payload. The transaction is instantly rejected by `VerifyTransaction`.

### 3. The Coinbase Exception
A **coinbase transaction** is a special type of system transaction used to reward the miner who successfully solves the Proof of Work puzzle (minting a fixed block subsidy reward of 10.0 tokens).
*   It is always the first transaction in each block and does not have a real sender (the `Sender` field is hardcoded as `"Coinbase"`).
*   Because the tokens are freshly minted by the system out of thin air and do not represent spending existing user funds, **coinbase transactions do not require a cryptographic signature** and skip the validation loop safely.

### 4. Immutable Block Packaging (SHA-256)
A block is the fundamental ledger unit that groups transactions together into sequential history.
*   **Industrial Data Integrity:** This implementation utilizes **SHA-256** (`crypto/sha256`), eliminating weak algorithms like MD5 which suffer from collision vulnerabilities.
*   **Full Serialization:** The block's identity hash is derived from its metadata headers *and* every single transaction and signature bundled inside it. Altering even a single character or amount inside a block after the fact instantly breaks its hash, propagating an invalid state that corrupts the chain downstream.
*   **Hex Encoding:** Raw binary hash bytes are securely formatted into clean hexadecimal strings to prevent runtime string-casting character corruption bugs.

### 5. Proof of Work (PoW) Consensus Puzzle
To securely commit blocks to the network, miners must solve a computational puzzle controlled by a variable named `Difficulty`.
*   **The Target:** A `big.Int` upper boundary calculated by left-shifting a bit representation. A higher difficulty sets a lower target threshold, requiring the final block hash to begin with a specific number of leading zeros.
*   **The Nonce:** An incrementing integer variable used by the miner to vary the data string payload. The mining engine (`proof.go`) repeatedly brute-forces hashes with different nonces until it hits a SHA-256 hash value that falls strictly below the target barrier.
*   **Consensus Audit:** Because hashing is slow but verification is instant, any network participant can call `Validate()` to verify that the miner performed the necessary work simply by passing the single successful nonce.

---

## 🗂️ Project Repository Architecture

```text
blockchain-project/
├── blockchain/
│   ├── block.go          # Block schemas, full transaction serialization, and SHA-256 hashing
│   ├── proof.go          # Big integer targets, Difficulty settings, and brute-force MineBlock loops
│   ├── blockchain.go     # Ledger state coordination, coinbase distribution, and consensus gates
│   └── wallet.go         # RSA key generation, base64 transaction signing and signature verification
├── go.mod
└── main.go               # Orchestration runner launcher
```

---

## 📄 License

This project is licensed under the **GNU Affero General Public License v3 (AGPL-3.0)**. 

### Reciprocity Requirements
Under the terms of the AGPL-3.0, if you modify this codebase or integrate it into a platform hosted over a network interface (including cloud APIs or private backend infrastructure), **you are legally bound to make your entire, modified server-side source code publicly available for free**. See the `LICENSE` file for full disclosures.
