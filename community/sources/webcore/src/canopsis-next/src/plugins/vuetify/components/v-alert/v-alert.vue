<script>
import { isCssColor } from 'vuetify/lib/util/colorUtils';
import VAlert from 'vuetify/lib/components/VAlert/VAlert';

import { isDarkColor } from '@/helpers/color';
import { getThemeVariableKeyFromCssVariable } from '@/helpers/entities/theme/entity';

export default {
  extends: VAlert,
  computed: {
    /**
     * Here we are adding background class adding based on the background color variable.
     */
    classes() {
      let backgroundColor = this.computedColor;

      if (isCssColor(backgroundColor)) {
        const backgroundColorKey = getThemeVariableKeyFromCssVariable(this.computedColor);

        backgroundColor = this.$vuetify.theme.currentTheme?.[backgroundColorKey] || '#fff';
      }

      const isDarkBackground = isDarkColor(backgroundColor);

      return {
        'v-sheet': true,
        'v-sheet--outlined': this.outlined,
        'v-sheet--shaped': this.shaped,
        ...this.themeClasses,
        ...this.elevationClasses,
        ...this.roundedClasses,
        'v-alert--dark-background': isDarkBackground,
        'v-alert--light-background': !isDarkBackground,
      };
    },

    /**
     * Instead of type we are returning css variable.
     */
    computedColor() {
      return this.color || `var(--v-${this.type}-background-base)`;
    },

    /**
     * Instead of condition we are returning css variable.
     */
    iconColor() {
      return `var(--v-${this.type}-base)`;
    },
  },
};
</script>

<style lang="scss">
.v-sheet {
  &.v-alert {
    --alert-color-dark: #212121;
    --alert-color-light: #d9d9d9;

    font-weight: 700;
    font-size: 14px;
    line-height: 22px;

    &--dark-background {
      color: var(--alert-color-light);
    }

    &--light-background {
      color: var(--alert-color-dark);
    }
  }
}
</style>
