<script>
import VCombobox from 'vuetify/es5/components/VCombobox';
import VSelect from 'vuetify/es5/components/VSelect/VSelect';

export default {
  extends: VCombobox,
  props: {
    blurOnCreate: {
      type: Boolean,
      default: false,
    },
  },
  methods: {
    onEnterDown(e) {
      e.preventDefault();

      VSelect.options.methods.onEnterDown.call(this);
      // If has menu index, let v-select-list handle
      if (this.getMenuIndex() > -1) {
        return;
      }

      this.updateSelf();

      /**
       * We've added this condition for the closing menu when we've pressed enter to create new one
       */
      if (this.blurOnCreate) {
        this.blur();
      }
    },
  },
};
</script>
