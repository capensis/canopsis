<script>
import VAutocomplete from 'vuetify/lib/components/VAutocomplete';

export default {
  extends: VAutocomplete,
  props: {
    clearOnEmptySearch: {
      type: Boolean,
      default: true,
    },
  },
  methods: {
    onInput(e) {
      if (this.selectedIndex > -1 || !e.target) return;
      const { target } = e;
      const { value } = target; // If typing and menu is not currently active

      if (target.value) this.activateMenu();
      if (!this.multiple && value === '' && !this.clearOnEmptySearch) this.deleteCurrentItem();
      this.internalSearch = value;
      this.badInput = target.validity && target.validity.badInput;
    },
  },
};
</script>
