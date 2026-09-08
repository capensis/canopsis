<template>
  <pattern-editor-field
    :patterns="patterns"
    :disabled="disabled"
    :readonly="readonly"
    :name="name"
    :type="PATTERN_TYPES.pbehavior"
    :required="required"
    :attributes="pbehaviorAttributes"
    :with-type="withType"
    :alarm-counter="alarmCounter"
    :entity-counter="entityCounter"
    class="c-pbehavior-patterns-field"
    @update:patterns="$emit('update:value', $event)"
  >
    <template #append-count="">
      <v-btn
        v-if="alarmCounter && alarmCounter.count"
        variant="text"
        small
        @click="showPatternAlarmsModal()"
      >
        {{ t('common.seeAlarms') }}
      </v-btn>
      <v-btn
        v-if="entityCounter && entityCounter.count"
        variant="text"
        small
        @click="showPatternEntitiesModal()"
      >
        {{ t('common.seeEntities') }}
      </v-btn>
    </template>
  </pattern-editor-field>
</template>

<script>
import { computed, onMounted, ref } from 'vue';

import {
  MAX_LIMIT,
  PATTERN_OPERATORS,
  PATTERN_TYPES,
  PBEHAVIOR_PATTERN_FIELDS,
  PBEHAVIOR_TYPE_TYPES,
} from '@/constants';

import { mergePatternAttributes } from '@/helpers/entities/pattern/fields/form';

import { useI18n } from '@/hooks/i18n';
import { usePbehavior } from '@/hooks/store/modules/pbehavior';

import PatternEditorField from '@/components/forms/fields/pattern/pattern-editor-field.vue';

export default {
  emits: ['update:patterns', 'update:value', 'show:alarms', 'show:entities'],
  components: { PatternEditorField },
  props: {
    patterns: {
      type: Object,
      required: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      required: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    withType: {
      type: Boolean,
      default: false,
    },
    checkCountName: {
      type: String,
      required: false,
    },
    readonly: {
      type: Boolean,
      default: false,
    },
    attributes: {
      type: Array,
      default: () => [],
    },
    alarmCounter: {
      type: Object,
      required: false,
    },
    entityCounter: {
      type: Object,
      required: false,
    },
  },
  setup(props, { emit }) {
    const { t } = useI18n();
    const {
      fetchPbehaviorsListWithoutStore,
      fetchPbehaviorReasonsListWithoutStore,
      fetchPbehaviorTypesListWithoutStore,
    } = usePbehavior();

    const pbehaviors = ref([]);
    const pbehaviorsPending = ref(false);

    const pbehaviorReasons = ref([]);
    const pbehaviorReasonsPending = ref(false);

    const pbehaviorTypes = ref([]);
    const pbehaviorTypesPending = ref(false);

    const nameOptions = computed(() => ({
      operators: [PATTERN_OPERATORS.equal, PATTERN_OPERATORS.notEqual],
      valueField: {
        is: 'c-select',
        props: {
          items: pbehaviors.value,
          loading: pbehaviorsPending.value,
          itemValue: '_id',
          itemTitle: 'name',
          ellipsis: true,
        },
      },
    }));

    const reasonOptions = computed(() => ({
      operators: [PATTERN_OPERATORS.equal, PATTERN_OPERATORS.notEqual],
      valueField: {
        is: 'c-select',
        props: {
          items: pbehaviorReasons.value,
          loading: pbehaviorReasonsPending.value,
          itemValue: '_id',
          itemTitle: 'name',
          ellipsis: true,
        },
      },
    }));

    const typeOptions = computed(() => ({
      operators: [PATTERN_OPERATORS.equal, PATTERN_OPERATORS.notEqual],
      valueField: {
        is: 'c-select',
        props: {
          items: pbehaviorTypes.value,
          loading: pbehaviorTypesPending.value,
          itemValue: '_id',
          itemTitle: 'name',
          ellipsis: true,
        },
      },
    }));

    const availablePbehaviorTypeTypes = computed(() => Object.values(PBEHAVIOR_TYPE_TYPES).map(type => ({
      value: type,
      text: t(`pbehavior.types.types.${type}`),
    })));

    const canonicalTypeOptions = computed(() => ({
      operators: [PATTERN_OPERATORS.equal, PATTERN_OPERATORS.notEqual],
      valueField: {
        is: 'c-select',
        props: {
          items: availablePbehaviorTypeTypes.value,
          itemTitle: 'text',
          itemValue: 'value',
        },
      },
    }));

    const defaultPbehaviorAttributes = computed(() => [
      {
        text: t('common.name'),
        value: PBEHAVIOR_PATTERN_FIELDS.id,
        options: nameOptions.value,
      },
      {
        text: t('common.reason'),
        value: PBEHAVIOR_PATTERN_FIELDS.reason,
        options: reasonOptions.value,
      },
      {
        text: t('common.type'),
        value: PBEHAVIOR_PATTERN_FIELDS.type,
        options: typeOptions.value,
      },
      {
        text: t('common.canonicalType'),
        value: PBEHAVIOR_PATTERN_FIELDS.canonicalType,
        options: canonicalTypeOptions.value,
      },
    ]);

    const pbehaviorAttributes = computed(() =>
      mergePatternAttributes(defaultPbehaviorAttributes.value, props.attributes));

    const showPatternAlarmsModal = () => {
      emit('show:alarms');
    };

    const showPatternEntitiesModal = () => {
      emit('show:entities');
    };

    const fetchPbehaviors = async () => {
      pbehaviorsPending.value = true;

      const { data: items } = await fetchPbehaviorsListWithoutStore({
        params: { limit: MAX_LIMIT },
      });

      pbehaviors.value = items;
      pbehaviorsPending.value = false;
    };

    const fetchPbehaviorReasons = async () => {
      pbehaviorReasonsPending.value = true;

      const { data: reasons } = await fetchPbehaviorReasonsListWithoutStore({
        params: { limit: MAX_LIMIT },
      });

      pbehaviorReasons.value = reasons;
      pbehaviorReasonsPending.value = false;
    };

    const fetchPbehaviorTypes = async () => {
      pbehaviorTypesPending.value = true;

      const { data: types } = await fetchPbehaviorTypesListWithoutStore({
        params: { limit: MAX_LIMIT },
      });

      pbehaviorTypes.value = types;
      pbehaviorTypesPending.value = false;
    };

    onMounted(() => {
      fetchPbehaviors();
      fetchPbehaviorReasons();
      fetchPbehaviorTypes();
    });

    return {
      t,
      PATTERN_TYPES,
      pbehaviorAttributes,
      showPatternAlarmsModal,
      showPatternEntitiesModal,
    };
  },
};
</script>
