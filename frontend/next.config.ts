import type { NextConfig } from "next";
import { config } from "dotenv";
import { resolve } from "node:path";

config({ path: resolve(process.cwd(), "../.env") });

const requiredEnv = ["NEXT_PUBLIC_API_URL"] as const;
const missing = requiredEnv.filter((key) => !process.env[key]);
if (missing.length > 0) {
  throw new Error(
    `Missing required env vars in root .env: ${missing.join(", ")}`
  );
}

const nextConfig: NextConfig = {
  /* config options here */
};

export default nextConfig;
