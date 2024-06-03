import {
  ref,
  provide,
  inject,
  onMounted,
  onBeforeUnmount,
} from 'vue';

import { AsyncBooting } from '@/services/async-booting';

/**
 * Hook for creating and managing an instance of AsyncBooting for asynchronous booting functionality.
 *
 * @returns {{ asyncBooting: Object }} An object containing the `asyncBooting` instance.
 */
export const useAsyncBootingParent = (itemsPerRender = 10, onFinishCallback) => {
  const asyncBooting = new AsyncBooting(itemsPerRender);

  provide('$asyncBooting', asyncBooting);

  onMounted(() => asyncBooting.run(itemsPerRender, onFinishCallback));
  onBeforeUnmount(() => asyncBooting.clear());

  return {
    asyncBooting,
  };
};

/**
 * Hook for handling async booting in a child component.
 *
 * @param {boolean} [bootedDefaultValue = false] - The default value for booted state.
 * @returns {{ booted: boolean, asyncBooting: Object }} An object containing `asyncBooting` and `booted` properties.
 */
export const useAsyncBootingChild = (bootedDefaultValue = false) => {
  const asyncBootingKey = Symbol('asyncBootingKey');
  const asyncBooting = inject('$asyncBooting');

  const booted = ref(bootedDefaultValue);

  const boot = () => booted.value = true;

  if (!bootedDefaultValue) {
    asyncBooting?.register(asyncBootingKey, boot);
  }

  onBeforeUnmount(() => asyncBooting?.unregister(asyncBootingKey));

  return {
    asyncBootingKey,
    asyncBooting,
    booted,
  };
};
