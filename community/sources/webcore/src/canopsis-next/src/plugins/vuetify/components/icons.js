import { camelCase } from 'lodash';

const requireModule = require.context(
  '!!svg-inline-loader?modules!../../../assets/images',
  true,
  /.*\.svg$/,
);
const icons = {};

requireModule.keys().forEach((fileName) => {
  const [, iconName] = fileName.match(/.+\/(.+).svg$/);

  icons[iconName] = {
    component: {
      name: camelCase(iconName),
      template: requireModule(fileName),
    },
  };
});

export default icons;
