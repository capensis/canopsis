<template lang="pug">
  v-btn(
    data-test="editTab",
    small,
    flat,
    icon,
    @click.prevent="showUpdateTabModal(tab)"
  )
    v-icon(small) edit
</template>

<script>
import { MODALS } from '@/constants';

export default {
  props: {
    tab: {
      type: Object,
      required: true,
    },
    updateTabMethod: {
      type: Function,
      default: () => {},
    },
  },

  methods: {
    showUpdateTabModal(tab) {
      this.$modals.show({
        name: MODALS.textFieldEditor,
        config: {
          title: this.$t('modals.viewTab.edit.title'),
          field: {
            name: 'text',
            label: this.$t('modals.viewTab.fields.title'),
            value: tab.title,
            validationRules: 'required',
          },
          action: (title) => {
            const newTab = { ...tab, title };

            return this.updateTabMethod(newTab);
          },
        },
      });
    },
  },
};
</script>
