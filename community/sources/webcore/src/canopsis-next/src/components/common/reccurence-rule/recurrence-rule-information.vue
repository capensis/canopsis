<template>
  <v-sheet color="grey lighten-2" class="c-word-break-all pa-3" rounded="xxl">
    <p>{{ recurrenceRuleString }}</p>
    <p>{{ recurrenceRuleText }}</p>
    <template v-if="exceptionItems.length > 0">
      <p>
        {{ $t('pbehavior.exdates.title') }}:
        <ul>
          <li v-for="item in exceptionItems" :key="item">
            {{ item }}
          </li>
        </ul>
      </p>
    </template>
  </v-sheet>
</template>

<script>
import { computed } from 'vue';
import { rrulestr } from 'rrule';

import { convertDateToString } from '@/helpers/date/date';

import { useI18n } from '@/hooks/i18n';

const MAX_EXCEPTION_ITEMS = 4;

export default {
  props: {
    rrule: {
      type: String,
      default: '',
    },
    exdates: {
      type: Array,
      default: () => [],
    },
    exceptions: {
      type: Array,
      default: () => [],
    },
  },
  setup(props) {
    const { t } = useI18n();

    const recurrenceRule = computed(() => {
      if (!props.rrule) {
        return null;
      }

      try {
        return rrulestr(props.rrule);
      } catch {
        return null;
      }
    });

    const recurrenceRuleString = computed(() => recurrenceRule.value?.toString() ?? '');
    const recurrenceRuleText = computed(() => recurrenceRule.value?.toText() ?? '');

    const exceptionItems = computed(() => {
      const result = props.exceptions.map(({ name }) => name);
      const preparedExdates = props.exdates.map(({ begin, end }) => `${convertDateToString(begin)} - ${convertDateToString(end)}`);

      result.push(...preparedExdates);

      if (result.length > MAX_EXCEPTION_ITEMS) {
        const count = result.length - MAX_EXCEPTION_ITEMS + 1;
        return result.slice(0, MAX_EXCEPTION_ITEMS - 1).concat([`+${count} ${t('common.more')}`]);
      }

      return result.length > MAX_EXCEPTION_ITEMS ? result.slice(0, MAX_EXCEPTION_ITEMS) : result;
    });

    return {
      recurrenceRule,
      recurrenceRuleString,
      recurrenceRuleText,
      exceptionItems,
    };
  },
};
</script>
