<template>
  <v-menu v-model="menu" max-height="300">
    <template #activator="{ props: activatorProps }">
      <v-btn
        v-bind="activatorProps"
        :loading="pending"
        :disabled="!availableExceptions.length && !hasMoreExceptions"
        color="primary"
        variant="outlined"
      >
        {{ $t('pbehavior.exceptions.choose') }}
      </v-btn>
    </template>

    <v-list density="comfortable">
      <v-list-item
        v-for="exception in availableExceptions"
        :key="exception._id"
        @click="addItemIntoArray(exception)"
      >
        <v-list-item-title>{{ exception.name }}</v-list-item-title>
      </v-list-item>

      <v-list-item v-if="pending">
        <v-list-item-title class="text-center">
          <v-progress-circular size="24" indeterminate />
        </v-list-item-title>
      </v-list-item>

      <div ref="appendElement" class="pbehavior-exceptions-list-menu__append" />
    </v-list>
  </v-menu>
</template>

<script>
import { ref, computed, onBeforeUnmount, onMounted, watch } from 'vue';

import { PAGINATION_LIMIT } from '@/config';

import { mapIds } from '@/helpers/array';

import { useArrayModelField } from '@/hooks/form/array-model-field';
import { usePbehaviorException } from '@/hooks/store/modules/pbehavior-exception';

export default {
  emits: ['update:value'],
  props: {
    value: {
      type: Array,
      required: true,
    },
  },
  setup(props, { emit }) {
    const { addItemIntoArray } = useArrayModelField(props, emit);
    const { fetchPbehaviorExceptionsListWithoutStore } = usePbehaviorException();

    const menu = ref(false);
    const exceptions = ref([]);
    const pending = ref(false);;
    const page = ref(1);
    const pageCount = ref(1);
    const appendElement = ref(null);
    let requestNumber = 0;

    const selectedExceptionsIds = computed(() => mapIds(props.value));

    const availableExceptions = computed(
      () => exceptions.value.filter(({ _id: id }) => !selectedExceptionsIds.value.includes(id)),
    );
    const hasMoreExceptions = computed(() => pageCount.value > page.value);

    const fetchList = async (nextPage = 1) => {
      const currentRequestNumber = ++requestNumber;

      try {
        pending.value = true;

        const { data = [], meta = {} } = await fetchPbehaviorExceptionsListWithoutStore({
          params: {
            page: nextPage,
            limit: PAGINATION_LIMIT,
          },
        });

        if (currentRequestNumber !== requestNumber) {
          return;
        }

        page.value = nextPage;
        pageCount.value = meta.page_count ?? 1;
        exceptions.value = nextPage === 1 ? data : [...exceptions.value, ...data];
      } catch (err) {
        console.error(err);
      } finally {
        if (currentRequestNumber === requestNumber) {
          pending.value = false;
        }
      }
    };

    const fetchMoreExceptions = () => {
      if (!pending.value && hasMoreExceptions.value) {
        fetchList(page.value + 1);
      }
    };

    const observer = new IntersectionObserver(([entry]) => {
      if (entry.isIntersecting && menu.value) {
        fetchMoreExceptions();
      }
    });

    watch(menu, (opened) => {
      if (opened) {
        fetchList();
      }
    });

    onMounted(fetchList);

    watch(appendElement, (element, oldElement) => {
      if (oldElement) {
        observer.unobserve(oldElement);
      }

      if (element) {
        observer.observe(element);
      }
    });

    onBeforeUnmount(() => observer.disconnect());

    return {
      menu,
      pending,
      availableExceptions,
      hasMoreExceptions,
      appendElement,
      addItemIntoArray,
    };
  },
};
</script>

<style lang="scss" scoped>
.pbehavior-exceptions-list-menu__append {
  height: 1px;
}
</style>
