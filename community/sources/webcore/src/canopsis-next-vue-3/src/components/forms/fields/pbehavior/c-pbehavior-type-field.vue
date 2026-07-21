<template>
  <c-autocomplete
    :model-value="modelValue"
    :items="items"
    :label="label || t('common.type')"
    :loading="pending"
    :name="name"
    :required="required"
    :disabled="disabled"
    :return-object="returnObject"
    :multiple="multiple"
    :chips="chips"
    :closable-chips="chips"
    :clearable="clearable"
    item-title="name"
    item-value="_id"
    @update:model-value="$emit('update:modelValue', $event)"
  />
</template>

<script>
import { onMounted, ref } from 'vue'

import { useI18n } from '@/hooks/i18n'
import { usePbehaviorType } from '@/hooks/store/modules/pbehavior-type'

export default {
  props: {
    modelValue: {
      type: [String, Object, Array],
      default: '',
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'type',
    },
    required: {
      type: Boolean,
      default: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    returnObject: {
      type: Boolean,
      default: false,
    },
    multiple: {
      type: Boolean,
      default: false,
    },
    chips: {
      type: Boolean,
      default: false,
    },
    clearable: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:modelValue'],

  setup() {
    const { t } = useI18n()
    const { fetchPbehaviorTypesFieldList } = usePbehaviorType()

    const items = ref([])
    const pending = ref(false)

    const fetchList = async () => {
      try {
        pending.value = true

        const { data = [] } = await fetchPbehaviorTypesFieldList({
          params: {
            paginate: false,
            with_hidden: true,
          },
        })

        items.value = data
      } finally {
        pending.value = false
      }
    }

    onMounted(fetchList)

    return {
      t,
      items,
      pending,
    }
  },
}
</script>
