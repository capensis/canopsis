<template lang="pug">
  widget-settings-item(:title="label")
    v-select(
      v-model="selectedTemplate",
      :items="availableTemplates",
      :label="$t('common.template')",
      :loading="pending"
    )
    c-columns-field(
      v-field="columns",
      :with-template="withTemplate",
      :with-html="withHtml",
      :with-color-indicator="withColorIndicator",
      :type="type",
      :alarm-infos="alarmInfos",
      :entity-infos="entityInfos",
      :infos-pending="infosPending"
    )
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { MAX_LIMIT } from '@/constants';

import WidgetSettingsItem from '@/components/sidebars/settings/partials/widget-settings-item.vue';

const { mapActions } = createNamespacedHelpers('view/widget/template');

export default {
  components: { WidgetSettingsItem },
  model: {
    prop: 'columns',
    event: 'input',
  },
  props: {
    label: {
      type: String,
      required: true,
    },
    type: {
      type: String,
      required: true,
    },
    columns: {
      type: [Array, Object],
      default: () => [],
    },
    template: {
      type: String,
      required: false,
    },
    withTemplate: {
      type: Boolean,
      default: false,
    },
    withHtml: {
      type: Boolean,
      default: false,
    },
    withColorIndicator: {
      type: Boolean,
      default: false,
    },
    alarmInfos: {
      type: Array,
      default: () => [],
    },
    entityInfos: {
      type: Array,
      default: () => [],
    },
    infosPending: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      widgetTemplates: [],
      pending: false,
    };
  },
  computed: {
    availableTemplates() {
      return [
        ...this.widgetTemplates.map(template => ({ value: template._id, text: template.title })),

        { value: 'custom', text: this.$t('common.custom') },
      ];
    },

    selectedTemplate: {
      get() {
        return this.template;
      },
      set(value) {
        this.$emit('update:widgetColumnsTemplate', value);
      },
    },
  },
  methods: {
    ...mapActions({
      fetchWidgetTemplatesListWithoutStore: 'fetchListWithoutStore',
    }),

    async fetchList() {
      this.pending = true;

      const { data } = await this.fetchWidgetTemplatesListWithoutStore({
        params: {
          type: this.type,
          limit: MAX_LIMIT,
        },
      });

      this.widgetTemplates = data;
      this.pending = false;
    },
  },
};
</script>
