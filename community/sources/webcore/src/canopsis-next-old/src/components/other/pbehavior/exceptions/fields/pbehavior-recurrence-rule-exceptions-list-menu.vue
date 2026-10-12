<template>
  <v-menu offset-y>
    <template #activator="{ on }">
      <v-btn
        :loading="pending"
        :disabled="!availableExceptions.length"
        color="primary"
        outlined
        v-on="on"
      >
        {{ $t('pbehavior.exceptions.choose') }}
      </v-btn>
    </template>
    <v-list dense>
      <v-list-item
        v-for="exception in availableExceptions"
        :key="exception._id"
        @click="addItemIntoArray(exception)"
      >
        <v-list-item-content>
          <v-list-item-title>{{ exception.name }}</v-list-item-title>
        </v-list-item-content>
      </v-list-item>
    </v-list>
  </v-menu>
</template>

<script>
import { ref, computed, onMounted } from 'vue';

import { MAX_LIMIT } from '@/constants';

import { mapIds } from '@/helpers/array';

import { useArrayModelField } from '@/hooks/form/array-model-field';
import { usePendingHandler } from '@/hooks/query/pending';
import { usePbehaviorException } from '@/hooks/store/modules/pbehavior-exception';

export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Array,
      required: true,
    },
  },
  setup(props, { emit }) {
    const { addItemIntoArray } = useArrayModelField(props, emit);
    const { fetchPbehaviorExceptionsListWithoutStore } = usePbehaviorException();

    const exceptions = ref([]);

    const selectedExceptionsIds = computed(() => mapIds(props.value));

    const availableExceptions = computed(() => exceptions.value.filter(
      ({ _id: id }) => !selectedExceptionsIds.value.includes(id),
    ));

    const { pending, handler: fetchList } = usePendingHandler(async () => {
      try {
        const { data } = await fetchPbehaviorExceptionsListWithoutStore({ params: { limit: MAX_LIMIT } });

        exceptions.value = data;
      } catch (err) {
        console.error(err);
      }
    });

    onMounted(fetchList);

    return {
      pending,
      availableExceptions,
      addItemIntoArray,
    };
  },
};
</script>
