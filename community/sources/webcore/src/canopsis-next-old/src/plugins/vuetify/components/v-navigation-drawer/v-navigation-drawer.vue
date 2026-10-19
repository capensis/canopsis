<script>
import VNavigationDrawer from 'vuetify/lib/components/VNavigationDrawer';

import ClickOutside from '../../directives/click-outside';
import overlayableMixin from '../../mixins/overlayable';

export default {
  directives: {
    ClickOutside,
  },
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
      const directives = [];

      if (!this.ignoreClickOutside) {
        directives.push({
          name: 'click-outside',
          value: {
            handler: () => {
              /**
               * We can't move call of customCloseConditional into closeConditional because here this method will call
               * more than once. But we need only once call
               */
              if (this.customCloseConditional && this.customCloseConditional() === false) {
                return;
              }

              this.isActive = false;
            },
            closeConditional: this.closeConditional,
            include: this.getOpenDependentElements,
          },
          modifiers: { same: true },
        });
      }

      if (!this.touchless && !this.stateless) {
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
