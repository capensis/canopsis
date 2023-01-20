<template lang="pug">
  v-layout(column)
    v-select(
      v-field="form.type",
      v-validate="'required'",
      :label="$t('common.type')",
      :items="availableTypes",
      :error-messages="errors.collect('type')",
      name="type"
    )
    v-text-field(
      v-field="form.title",
      v-validate="'required'",
      :label="$t('common.name')",
      :error-messages="errors.collect('title')",
      name="title"
    )
    span.body-2.my-2 {{ $tc('common.column', 2) }}
    c-columns-field(
      v-field="form.columns",
      :type="form.type",
      :alarm-infos="alarmInfos",
      :entity-infos="entityInfos",
      :infos-pending="infosPending",
      with-color-indicator,
      with-template,
      with-html
    )
</template>

<script>
import { ENTITIES_TYPES } from '@/constants';

import { widgetColumnsInfosMixin } from '@/mixins/widget/columns/infos';

export default {
  inject: ['$validator'],
  mixins: [widgetColumnsInfosMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    availableTypes() {
      return [ENTITIES_TYPES.alarm, ENTITIES_TYPES.entity].map(value => ({
        value: ENTITIES_TYPES.alarm,
        text: this.$t(`entities.${value}`),
      }));
    },
  },
};
</script>
