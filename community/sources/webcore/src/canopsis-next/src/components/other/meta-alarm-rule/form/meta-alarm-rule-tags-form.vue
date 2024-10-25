<template>
  <v-layout class="gap-3" column>
    <span class="text-subtitle-1 font-weight-bold">{{ $tc('common.tag', 2) }}</span>
    <v-layout>
      <c-enabled-field
        v-field="form.copy_from_children"
        :label="$t('metaAlarmRule.copyTagsFromChildren')"
        hide-details
      />
    </v-layout>
    <v-layout>
      <c-enabled-field
        v-model="filterByLabelEnabled"
        :label="$t('metaAlarmRule.filterByLabelEnabled')"
        class="mr-4 pt-4"
      >
        <template #append>
          <c-help-icon
            :text="$t('metaAlarmRule.filterByLabelEnabledTooltip')"
            icon="help"
            top
          />
        </template>
      </c-enabled-field>
      <v-fade-transition>
        <v-combobox
          v-if="filterByLabelEnabled"
          v-field="form.filter_by_label"
          v-validate="'required'"
          :label="$t('common.label')"
          :error-messages="errors.collect('filter_by_label')"
          class="mt-0"
          name="filter_by_label"
          chips
          multiple
        />
      </v-fade-transition>
    </v-layout>
  </v-layout>
</template>

<script>
import { ref } from 'vue';

export default {
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
    variables: {
      type: Array,
      default: () => [],
    },
  },
  setup(props) {
    const filterByLabelEnabled = ref(!!props.form.filter_by_label?.length);

    return {
      filterByLabelEnabled,
    };
  },
};
</script>
