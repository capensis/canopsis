import { DEFAULT_SANITIZE_OPTIONS } from '@/config';

export const TEXT_WIDGET_SANITIZE_OPTIONS = {
  ...DEFAULT_SANITIZE_OPTIONS,

  allowedTags: [...DEFAULT_SANITIZE_OPTIONS.allowedTags, 'iframe'],
  allowedAttributes: {
    ...DEFAULT_SANITIZE_OPTIONS.allowedAttributes,

    iframe: [
      'allowtransparency',
      'frameborder',
      'hspace',
      'marginheight',
      'marginwidth',
      'sandbox',
      'scrolling',
      'seamless',
      'src',
      'srcdoc',
      'vspace',
    ],
  },
};
