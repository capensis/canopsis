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

/**
 * Icons available as Material Symbol ligatures (no SVG needed).
 * These were previously custom SVGs but now use the font directly.
 */
const MATERIAL_SYMBOL_ICONS = [
  'alt_route',
  'bookmark_add',
  'bookmark_remove',
  'build_circle',
  'collapse_all',
  'density_large',
  'density_medium',
  'density_small',
  'engineering',
  'expand_all',
  'mark_unread_chat_alt',
  'mediation',
  'motion_photos_paused',
  'numbers',
  'published_with_changes',
  'restart_alt',
  'sticky_note_2',
  'storage',
  'variables',
  'webhook',
];

MATERIAL_SYMBOL_ICONS.forEach((name) => {
  icons[name] = name;
});

export default icons;
