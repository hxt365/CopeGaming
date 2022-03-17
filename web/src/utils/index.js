import { Buffer } from "buffer";

export const decodeBase64 = (b64str) => {
  const b = Buffer.from(b64str, "base64");
  return b.toString();
};

export const encodeBase64 = (str) => {
  const b = Buffer.from(str);
  return b.toString("base64");
};
