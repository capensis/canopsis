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

const getIconNameFromPath = path => path.match(/\/([^/]+)\.svg$/)[1];

const registerAlias = (name, value) => {
  customIconAliases[name] = value;
  customIconAliases[`vuetify.icons.${name}`] = value;
};

export const registerCustomIcon = (name, template, meta = {}) => {
  customIconsRegistry[name] = {
    meta,
    component: {
      template,
    },
  };

  registerAlias(name, `canopsis:${name}`);
};

export const unregisterCustomIcon = (name) => {
  delete customIconsRegistry[name];
  delete customIconAliases[name];
  delete customIconAliases[`vuetify.icons.${name}`];
};

Object.entries(staticIconModules).forEach(([path, template]) => {
  registerCustomIcon(getIconNameFromPath(path), template);
});

MATERIAL_SYMBOL_ICONS.forEach((name) => {
  registerAlias(name, name);
});

const getIconProps = (props) => {
  const iconProps = { ...props };

  delete iconProps.tag;
  delete iconProps.icon;

  return iconProps;
};

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
