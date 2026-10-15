import { h } from 'vue';
import { aliases as mdAliases, md } from 'vuetify/iconsets/md';

const OUTLINED_ICON_SUFFIX = '_outline';

const staticIconModules = import.meta.glob('@/assets/images/*.svg', {
  eager: true,
  query: '?raw',
  import: 'default',
});

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

export const customIconsRegistry = {};

export const customIconAliases = {
  ...mdAliases,
};

let liveAliases = customIconAliases;

/**
 * Extracts the icon name from an SVG asset path.
 *
 * @param {string} path Absolute or aliased path ending with `.svg`.
 * @returns {string} File name without extension.
 */
const getIconNameFromPath = path => path.match(/\/([^/]+)\.svg$/)[1];

/**
 * Writes an alias into both the local map and Vuetify's live aliases object.
 *
 * @param {string} key Alias key.
 * @param {string} value Icon set identifier or Material icon name.
 */
const setAlias = (key, value) => {
  customIconAliases[key] = value;
  liveAliases[key] = value;
};

/**
 * Removes an alias from both the local map and Vuetify's live aliases object.
 *
 * @param {string} key Alias key to remove.
 */
const deleteAlias = (key) => {
  delete customIconAliases[key];
  delete liveAliases[key];
};

/**
 * Registers an icon alias and the Vuetify-prefixed variants used by templates.
 *
 * @param {string} name Icon name.
 * @param {string} value Icon set identifier or Material icon name.
 */
const registerAlias = (name, value) => {
  setAlias(name, value);
  setAlias(`vuetify.icons.${name}`, value);
  setAlias(`vuetify.icon.${name}`, value);
};

/**
 * Points alias writes at Vuetify's live aliases object.
 * createIcons copies aliases at init, so runtime uploads must update that copy.
 *
 * @param {Object<string, string>} [aliases] Live aliases map from Vuetify.
 */
export const bindIconAliases = (aliases) => {
  if (!aliases) {
    return;
  }

  liveAliases = aliases;
  Object.assign(liveAliases, customIconAliases);
};

/**
 * Registers a custom SVG icon and aliases it to the `canopsis` icon set.
 *
 * @param {string} name Icon name.
 * @param {string} template Raw SVG markup.
 * @param {Object} [meta={}] Extra render flags.
 * @param {boolean} [meta.fillBorder] Whether the icon should use a fill-border class.
 */
export const registerCustomIcon = (name, template, meta = {}) => {
  customIconsRegistry[name] = {
    meta,
    component: {
      template,
    },
  };

  registerAlias(name, `canopsis:${name}`);
};

/**
 * Removes a custom icon from the registry and deletes its aliases.
 *
 * @param {string} name Icon name to unregister.
 */
export const unregisterCustomIcon = (name) => {
  delete customIconsRegistry[name];
  deleteAlias(name);
  deleteAlias(`vuetify.icons.${name}`);
  deleteAlias(`vuetify.icon.${name}`);
};

Object.entries(staticIconModules).forEach(([path, template]) => {
  registerCustomIcon(getIconNameFromPath(path), template);
});

MATERIAL_SYMBOL_ICONS.forEach((name) => {
  registerAlias(name, name);
});

/**
 * Strips Vuetify icon renderer props that must not be forwarded to the DOM node.
 *
 * @param {Object} props Icon renderer props.
 * @param {string} props.tag HTML tag used for the icon wrapper.
 * @param {string} props.icon Icon name to render.
 * @returns {Object} Props without `tag` and `icon`.
 */
const getIconProps = (props) => {
  const iconProps = { ...props };

  delete iconProps.tag;
  delete iconProps.icon;

  return iconProps;
};

/**
 * Renders a registered custom SVG icon, or falls back to the icon name as text.
 *
 * @param {Object} props Icon renderer props.
 * @param {string} props.tag HTML tag used for the icon wrapper.
 * @param {string} props.icon Icon name to look up in `customIconsRegistry`.
 * @returns {VNode} Rendered icon node.
 */
const CanopsisIcon = (props) => {
  const icon = customIconsRegistry[props.icon];
  const iconProps = getIconProps(props);

  if (!icon) {
    return h(props.tag, iconProps, props.icon);
  }

  return h(props.tag, {
    ...iconProps,
    class: [
      iconProps.class,
      {
        'v-icon--fill-border': icon.meta?.fillBorder,
        'v-icon--no-fill': props.icon.endsWith(OUTLINED_ICON_SUFFIX),
      },
    ],
    innerHTML: icon.component.template,
  });
};

export const canopsisIcons = {
  component: CanopsisIcon,
};

export const iconSets = {
  md,
  canopsis: canopsisIcons,
};
