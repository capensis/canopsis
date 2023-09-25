<script>
import VNavigationDrawer from 'vuetify/es5/components/VNavigationDrawer';

import overlayableMixin from '../../mixins/overlayable';

export default {
  extends: VNavigationDrawer,
  mixins: [overlayableMixin],
  props: {
    ignoreClickOutside: {
      type: Boolean,
      default: false,
    },
    customCloseConditional: {
      type: Function,
      default: null,
    },
  },
  computed: {
    reactsToClick() {
      return !this.stateless && !this.permanent && (this.isMobile || this.temporary) && !this.ignoreClickOutside;
    },
  },
  methods: {
    genDirectives() {
      const directives = [{
        name: 'click-outside',
        value: () => {
          /**
           * We can't move call of customCloseConditional into closeConditional because here this method will call
           * more than once. But we need only once call
           */
          if (this.customCloseConditional && this.customCloseConditional() === false) {
            return;
          }

          this.isActive = false;
        },
        args: {
          closeConditional: this.closeConditional,
          include: this.getOpenDependentElements,
        },
      }];

      if (!this.touchless) {
        directives.push({
          name: 'touch',
          value: {
            parent: true,
            left: this.swipeLeft,
            right: this.swipeRight,
          },
        });
      }

      return directives;
    },
  },
};
</script>
