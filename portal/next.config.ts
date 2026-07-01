import type { NextConfig } from "next";
import path from "path";

const nextConfig: NextConfig = {
  transpilePackages: ["@douyinfe/semi-ui", "@douyinfe/semi-icons"],
  turbopack: {
    root: path.join(process.cwd()),
  },
};

export default nextConfig;
