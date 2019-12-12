<template lang="pug">
  div
    v-layout(wrap, justify-center)
      v-flex(xs11)
        v-text-field(
          v-field="form.name",
          v-validate="'required'",
          :label="$t('common.name')",
          :error-messages="errors.collect('name')",
          name="name",
          data-test="viewFieldName"
        )
        v-text-field(
          v-field="form.title",
          v-validate="'required'",
          :label="$t('common.title')",
          :error-messages="errors.collect('title')",
          name="title",
          data-test="viewFieldTitle"
        )
        v-text-field(
          v-field="form.description",
          :label="$t('common.description')",
          name="description",
          data-test="viewFieldDescription"
        )
        v-switch(
          v-field="form.enabled",
          :label="$t('common.enabled')",
          data-test="viewFieldEnabled"
        )
    v-layout(wrap, justify-center)
      v-flex(xs11)
        v-combobox(
          v-field="form.tags",
          :label="$t('modals.view.fields.groupTags')",
          data-test="viewFieldGroupTags",
          tags,
          clearable,
          multiple,
          append-icon,
          chips,
          deletable-chips
        )
        v-combobox(
          ref="combobox",
          v-validate="'required'",
          :value="groupName",
          :items="groupNames",
          :label="$t('modals.view.fields.groupIds')",
          :error-messages="errors.collect('group')",
          name="group",
          data-test="viewFieldGroupId",
          @change="changeGroupName"
        )
          template(slot="no-data")
            v-list-tile
              v-list-tile-content
                v-list-tile-title(v-html="$t('modals.view.noData')")
</template>

<script>

import vuetifyComboboxMixin from '@/mixins/vuetify/combobox';

export default {
  inject: ['$validator'],
  mixins: [vuetifyComboboxMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    groupName: {
      type: String,
      default: '',
    },
    groups: {
      type: Array,
      default: () => [],
    },
  },
  computed: {
    groupNames() {
      return this.groups.map(group => group.name);
    },
  },
  methods: {
    changeGroupName(value) {
      this.$emit('update:groupName', value);
      this.closeComboboxMenuOnChange();
    },
  },
};
</script>
