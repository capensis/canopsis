const requireModule = require.context('.', true, /.*\.svg$/);
const assets = [];

requireModule.keys().forEach((fileName) => {
  assets.push(requireModule(fileName));
});

export default assets;
