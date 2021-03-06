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
        enabled-field(v-field="form.enabled", data-test="viewFieldEnabled")
        periodic-refresh-field(v-model="form.periodicRefresh", :label="$t('modals.view.fields.periodicRefresh')")
    v-layout(wrap, justify-center)
      v-flex(xs11)
        v-combobox(
          v-field="form.tags",
          :label="$t('modals.view.fields.groupTags')",
          data-test="viewFieldGroupTags",
          append-icon="",
          tags,
          clearable,
          multiple,
          chips,
          deletable-chips
        )
        v-combobox(
          v-validate="'required'",
          :value="groupName",
          :items="groupNames",
          :label="$t('modals.view.fields.groupIds')",
          :error-messages="errors.collect('group')",
          name="group",
          data-test="viewFieldGroupId",
          blur-on-create,
          @change="changeGroupName"
        )
          template(slot="no-data")
            v-list-tile
              v-list-tile-content
                v-list-tile-title(v-html="$t('modals.view.noData')")
</template>

<script>
import PeriodicRefreshField from '@/components/forms/fields/periodic-refresh-field.vue';
import EnabledField from '@/components/forms/fields/enabled-field.vue';

export default {
  components: { EnabledField, PeriodicRefreshField },
  inject: ['$validator'],
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
    },
  },
};
</script>
