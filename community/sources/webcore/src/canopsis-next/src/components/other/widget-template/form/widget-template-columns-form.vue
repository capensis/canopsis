<template lang="pug">
  v-layout(column)
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
      :type="entityType",
      :alarm-infos="alarmInfos",
      :entity-infos="entityInfos",
      :infos-pending="infosPending",
      with-color-indicator,
      with-template,
      with-html
    )
</template>

<script>
import { ENTITIES_TYPES, WIDGET_TEMPLATES_TYPES } from '@/constants';

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
  computed: {
    entityType() {
      return this.form.type === WIDGET_TEMPLATES_TYPES.alarmColumns
        ? ENTITIES_TYPES.alarm
        : ENTITIES_TYPES.entity;
    },
  },
};
</script>
