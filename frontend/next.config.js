
module.exports = {
  async rewrites() {
    return [
      {
        source: '/api/:path*',      
        destination: `${process.env.BACKEND_URL}/:path*`,      
      },
    ];
  },
};
/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'export',
}
module.exports = nextConfig