<template>
  <v-layout column>
    <c-label>{{ $tc('common.tag', 2) }}</c-label>
    <v-layout>
      <c-enabled-field
        v-field="form.copy_from_children"
        :label="$t('metaAlarmRule.copyTagsFromChildren')"
        hide-details
      />
    </v-layout>
    <v-layout v-if="form.copy_from_children" class="gap-3" align-start>
      <c-enabled-field
        v-model="filterByLabelEnabled"
        :label="$t('metaAlarmRule.filterByLabelEnabled')"
        hide-details
      >
        <template #append>
          <c-help-icon
            :text="$t('metaAlarmRule.filterByLabelEnabledTooltip')"
            icon="help"
            top
          />
        </template>
      </c-enabled-field>
      <v-expand-transition>
        <c-alarm-tag-field
          v-if="filterByLabelEnabled"
          v-field="form.filter_by_label"
          :label="$t('common.label')"
          name="filter_by_label"
          multiple
          addable
          combobox
          required
          hide-selected
          only-labels
        >
          <template #no-data="">
            <v-list-item>
              <v-list-item-content>
                <v-list-item-title v-html="$t('common.pressEnterToApply')" />
              </v-list-item-content>
            </v-list-item>
          </template>
        </c-alarm-tag-field>
      </v-expand-transition>
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
