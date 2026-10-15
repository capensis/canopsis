const requireModule = require.context(
  '!!svg-inline-loader?modules!.',
  true,
  /.*\.svg$/,
);
const assetGroups = {};

requireModule.keys().forEach((fileName) => {
  const [, groupName] = fileName.match(/\.\/(.+)\/.+.svg$/);
  const file = requireModule(fileName);

  if (assetGroups[groupName]) {
    assetGroups[groupName].push(file);
  } else {
    assetGroups[groupName] = [file];
  }
});

export default assetGroups;
