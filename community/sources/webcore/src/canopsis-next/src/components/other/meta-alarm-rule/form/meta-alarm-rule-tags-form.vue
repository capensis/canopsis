<template>
  <v-layout class="gap-3" column>
    <span class="text-subtitle-1 font-weight-bold">{{ $tc('common.tag', 2) }}</span>
    <v-layout>
      <c-enabled-field
        v-field="form.copy_from_children"
        :label="$t('metaAlarmRule.copyTagsFromChildren')"
        hide-details
        with-background
      />
    </v-layout>
    <v-layout v-if="form.copy_from_children">
      <c-enabled-field
        v-model="filterByLabelEnabled"
        :label="$t('metaAlarmRule.filterByLabelEnabled')"
        class="mr-4 pt-4"
        with-background
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
