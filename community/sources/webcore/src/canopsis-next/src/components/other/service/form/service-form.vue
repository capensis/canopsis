<template>
  <v-layout column>
    <c-name-field
      v-field="form.name"
      required
    />
    <v-layout>
      <v-flex
        class="pr-3"
        xs6
      >
        <c-entity-category-field
          v-field="form.category"
          class="mt-1"
          addable
          required
        />
      </v-flex>
      <v-flex
        class="pr-3"
        xs4
      >
        <c-alarm-state-field
          v-field="form.sli_avail_state"
          :label="$t('service.availabilityState')"
          required
        />
      </v-flex>
      <v-flex xs2>
        <c-impact-level-field
          v-field="form.impact_level"
          required
        />
      </v-flex>
    </v-layout>
    <c-coordinates-field
      v-field="form.coordinates"
      row
    />
    <text-editor-field
      v-validate="'required'"
      v-field="form.output_template"
      :label="$t('service.outputTemplate')"
      :error-messages="errors.collect('output_template')"
      :variables="outputVariables"
      name="output_template"
    />
    <c-enabled-field v-field="form.enabled" />
    <entity-state-setting
      :form="form"
      :preparer="prepareStateSettingForm"
    />
    <v-tabs
      slider-color="primary"
      centered
    >
      <v-tab :class="{ 'error--text': errors.has('entity_patterns') }">
        {{ $t('common.entityPatterns') }}
      </v-tab>
      <v-tab-item>
        <c-patterns-field
          v-field="form.patterns"
          :entity-attributes="entityAttributes"
          class="mt-2"
          with-entity
          entity-counters-type
        />
      </v-tab-item>
      <v-tab
        :disabled="advancedJsonWasChanged"
        class="validation-header"
      >
        {{ $t('entity.manageInfos') }}
      </v-tab>
      <v-tab-item>
        <manage-infos v-field="form.infos" />
      </v-tab-item>
    </v-tabs>
  </v-layout>
</template>

<script>
import { get } from 'lodash';

import {
  ENTITY_PATTERN_FIELDS,
  SERVICE_WEATHER_STATE_COUNTERS,
  SERVICE_WEATHER_TEMPLATE_COUNTERS_BY_STATE_COUNTERS,
} from '@/constants';

import ManageInfos from '@/components/widgets/context/manage-infos.vue';
import TextEditorField from '@/components/forms/fields/text-editor-field.vue';
import EntityStateSetting from '@/components/other/state-setting/entity-state-setting.vue';

export default {
  inject: ['$validator'],
  components: {
    TextEditorField,
    ManageInfos,
    EntityStateSetting,
  },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    prepareStateSettingForm: {
      type: Function,
      default: data => data,
    },
  },
  computed: {
    outputVariables() {
      const messages = this.$t('serviceWeather.stateCounters');

      return Object.values(SERVICE_WEATHER_STATE_COUNTERS).map(field => ({
        text: messages[field],
        value: SERVICE_WEATHER_TEMPLATE_COUNTERS_BY_STATE_COUNTERS[field],
      }));
    },

    advancedJsonWasChanged() {
      return get(this.fields, ['advancedJson', 'changed']);
    },

    entityAttributes() {
      return [
        {
          value: ENTITY_PATTERN_FIELDS.lastEventDate,
          options: { disabled: true },
        },
        {
          value: ENTITY_PATTERN_FIELDS.connector,
          options: { disabled: true },
        },
        {
          value: ENTITY_PATTERN_FIELDS.componentInfos,
          options: { disabled: true },
        },
      ];
    },
  },
};
</script>
