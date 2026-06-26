import { DEFAULT_SANITIZE_OPTIONS } from '@/config';

export const TEXT_WIDGET_SANITIZE_OPTIONS = {
  ...DEFAULT_SANITIZE_OPTIONS,

  ADD_TAGS: [
    ...DEFAULT_SANITIZE_OPTIONS.ADD_TAGS,
    'iframe',
  ],

  ADD_ATTR: [
    ...DEFAULT_SANITIZE_OPTIONS.ADD_ATTR,
    'src',
    'width',
    'height',
    'allow',
    'allowfullscreen',
    'frameborder',
  ],
};
