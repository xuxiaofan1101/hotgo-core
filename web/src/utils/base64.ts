import CryptoJS from 'crypto-js';

export const base64 = {
  // base64编码
  encode(word: string): string {
    const src = CryptoJS.enc.Utf8.parse(word);
    return CryptoJS.enc.Base64.stringify(src);
  },
  // base64解码
  decode(word: string): string {
    const src = CryptoJS.enc.Base64.parse(word);
    return CryptoJS.enc.Utf8.stringify(src);
  },
  // base64URl编码
  urlEncode(word: string): string {
    const src = CryptoJS.enc.Utf8.parse(word);
    const b64 = CryptoJS.enc.Base64.stringify(src);
    return b64.replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/g, '');
  },
  // base64URL解码
  urlDecode(word: string): string {
    let b64 = word.replace(/-/g, '+').replace(/_/g, '/');
    while (b64.length % 4 !== 0) {
      b64 += '=';
    }
    const src = CryptoJS.enc.Base64.parse(b64);
    return CryptoJS.enc.Utf8.stringify(src);
  },
};
