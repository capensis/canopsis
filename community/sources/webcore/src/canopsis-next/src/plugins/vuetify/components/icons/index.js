const requireModule = require.context(
  '.',
  true,
  /.*\.vue$/,
);
const icons = {};

requireModule.keys().forEach((fileName) => {
  const [, iconName] = fileName.match(/.+\/(.+).vue$/);

  icons[iconName] = {
    component: requireModule(fileName).default,
  };
});

export default icons;
