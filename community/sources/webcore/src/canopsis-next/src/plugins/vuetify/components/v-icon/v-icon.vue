<script>
import { remapInternalIcon } from 'vuetify/lib/util/helpers';
import VIcon from 'vuetify/lib/components/VIcon/VIcon';

/**
 * We added functionality to add custom class `v-icon--fill-border` by meta.fillBorder flag
 * This class is using in `v-chip` component
 */
export default {
  extends: VIcon,
  render(h, config) {
    const iconName = config.data?.domProps?.textContent
      || config.data?.domProps?.innerHTML
      || config.children[0]?.text?.trim?.()
      || '';

    const icon = remapInternalIcon(config.parent, iconName);

    config.data.class = {
      ...config.data.class,

      'v-icon--fill-border': icon?.meta?.fillBorder,
    };

    return VIcon.options.render(h, config);
  },
};
</script>
