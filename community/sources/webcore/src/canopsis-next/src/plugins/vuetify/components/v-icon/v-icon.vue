<script>
import { isString } from 'lodash';
import { remapInternalIcon } from 'vuetify/lib/util/helpers';
import VIcon from 'vuetify/lib/components/VIcon/VIcon';

const OUTLINED_ICON_SUFFIX = '_outline';

/**
 * We added functionality to add custom class `v-icon--fill-border` by meta.fillBorder flag
 * This class is using in `v-chip` component
 *
 * We also detect icons with `_outline` suffix and apply `v-icon--no-fill`
 * to render them in outlined style (FILL 0) while the rest use filled style (FILL 1).
 */
export default {
  extends: VIcon,
  render(h, config) {
    const iconName = config.data?.domProps?.textContent
      || config.data?.domProps?.innerHTML
      || config.children[0]?.text?.trim?.()
      || '';

    const icon = remapInternalIcon(config.parent, iconName);
    const resolvedName = isString(icon) ? icon : iconName;

    if (isString(config.data.class)) {
      config.data.class = {
        [config.data.class]: true,
      };
    }

    config.data.class = {
      ...config.data.class,

      'v-icon--fill-border': icon?.meta?.fillBorder,
      'v-icon--no-fill': resolvedName.endsWith(OUTLINED_ICON_SUFFIX),
    };

    return VIcon.options.render(h, config);
  },
};
</script>
