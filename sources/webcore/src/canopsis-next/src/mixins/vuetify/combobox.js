export default {
  methods: {
    closeComboboxMenuOnChange(comboboxRefKey = 'combobox') {
      if (this.$refs[comboboxRefKey]) {
        this.$nextTick(() => this.$refs[comboboxRefKey].$refs.menu.isActive = false);
      }
    },
  },
};
