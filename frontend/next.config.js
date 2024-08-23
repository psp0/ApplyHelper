// module.exports = {
//   async rewrites() {
//     return [
//       {
//         source: '/api/:path*',      
//         destination: `${process.env.BACKEND_URL}/:path*`,      
//       },
//     ];
//   },
// };
/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: false,
  output: 'export',
  env: {
    BACKEND_URL: process.env.NODE_ENV === 'development' 
      ? process.env.DEV_BACKEND_URL 
      : process.env.PROD_BACKEND_URL,
  },
}
module.exports = nextConfig