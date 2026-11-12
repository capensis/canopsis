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

/**
 * Only the exact expressions emitted by built-in helpers may reach Vue's compiler.
 */
export const TEMPLATE_VUE_ATTRIBUTES = {
  'c-alarm-tags-chips': {
    ':alarm': 'alarm',
    ':selected-tags': 'selectedTags',
    '@select': "$emit('select:tag', $event)",
    '@close': "$emit('remove:tag', $event)",
  },
  'c-entity-tags-chips': { ':entity': 'entity' },
  'c-links-list': { 'v-if': 'links && hasAccessForLinks', ':links': 'links' },
  'service-entities-list': {
    ':service': 'service',
    ':service-entities': 'serviceEntities',
    ':widget-parameters': 'widgetParameters',
    ':options': 'options',
    ':total-items': 'totalItems',
    ':actions-requests': 'actionsRequests',
    '@refresh': 'refreshEntities',
    '@update:options': 'updateOptions',
    '@add:action': 'addAction',
  },
  'c-remediation-instruction-execution-see-details': { ':execution': 'execution' },
};

/**
 * SVG has a separate, case-sensitive policy. External resources and animation
 * elements are deliberately excluded; paint servers and use references are local.
 */
export const CUSTOM_ICON_SANITIZE_OPTIONS = {
  allowedTags: [
    'svg', 'g', 'path', 'rect', 'circle', 'ellipse', 'line', 'polyline', 'polygon',
    'defs', 'linearGradient', 'radialGradient', 'stop', 'clipPath', 'mask',
    'use', 'symbol', 'title', 'desc', 'text', 'tspan', 'filter', 'feFlood', 'feColorMatrix',
    'feOffset', 'feGaussianBlur', 'feComposite', 'feBlend', 'feMerge', 'feMergeNode',
  ],
  allowedAttributes: {
    '*': [
      'id', 'class', 'xmlns', 'xmlns:xlink', 'viewBox', 'width', 'height', 'x', 'y', 'x1', 'y1', 'x2', 'y2',
      'cx', 'cy', 'r', 'rx', 'ry', 'd', 'points', 'transform', 'preserveAspectRatio',
      'fill', 'fill-rule', 'fill-opacity', 'stroke', 'stroke-width', 'stroke-linecap',
      'stroke-linejoin', 'stroke-miterlimit', 'stroke-dasharray', 'stroke-dashoffset', 'stroke-opacity',
      'opacity', 'clip-path', 'clip-rule', 'mask', 'maskUnits', 'maskContentUnits',
      'gradientUnits', 'gradientTransform', 'spreadMethod', 'offset', 'stop-color', 'stop-opacity',
      'fx', 'fy', 'fr', 'clipPathUnits', 'href', 'xlink:href', 'vector-effect', 'filter',
      'filterUnits', 'primitiveUnits', 'color-interpolation-filters', 'result', 'in', 'in2',
      'type', 'values', 'stdDeviation', 'operator', 'k1', 'k2', 'k3', 'k4', 'mode', 'flood-color', 'flood-opacity',
      'font-family', 'font-size', 'font-weight', 'text-anchor', 'dominant-baseline', 'dx', 'dy',
    ],
  },
  parser: { lowerCaseTags: false, lowerCaseAttributeNames: false, decodeEntities: true },
  disallowedTagsMode: 'discard',
  nonTextTags: ['script', 'style', 'foreignObject', 'iframe', 'object', 'embed', 'animate', 'set'],
  allowedSchemes: [],
};
