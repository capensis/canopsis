<template>
  <v-layout column>
    <c-name-field
      v-field="form.title"
      :label="$t('common.title')"
      name="title"
      required
    />
    <v-text-field
      v-field="form.description"
      :label="$t('common.description')"
      name="description"
    />
    <c-enabled-field v-field="form.enabled" />
    <periodic-refresh-field
      v-field="form.periodic_refresh"
      :label="$t('modals.view.fields.periodicRefresh')"
    />
    <v-combobox
      v-field="form.tags"
      :label="$t('modals.view.fields.groupTags')"
      append-icon=""
      tags
      clearable
      multiple
      chips
      deletable-chips
    />
    <v-combobox
      v-field="form.group"
      v-validate="'required'"
      :items="groups"
      :label="$t('modals.view.fields.groupIds')"
      :error-messages="errors.collect('group')"
      item-text="title"
      item-value="_id"
      name="group"
      return-object
      blur-on-create
    >
      <template #no-data="">
        <v-list-item>
          <v-list-item-content>
            <v-list-item-title v-html="$t('modals.view.noData')" />
          </v-list-item-content>
        </v-list-item>
      </template>
    </v-combobox>
  </v-layout>
</template>

<script>
import PeriodicRefreshField from '@/components/forms/fields/periodic-refresh-field.vue';

export default {
  inject: ['$validator'],
  components: { PeriodicRefreshField },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    groups: {
      type: Array,
      default: () => [],
    },
  },
};
</script>
