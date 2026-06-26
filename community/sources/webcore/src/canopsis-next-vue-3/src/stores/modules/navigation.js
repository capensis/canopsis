import { ref } from 'vue';
import { defineStore } from 'pinia';

export const useNavigationStore = defineStore('navigation', () => {
  const isEditingMode = ref(false);

  const toggleEditingMode = () => {
    isEditingMode.value = !isEditingMode.value;
  };

  return {
    isEditingMode,
    toggleEditingMode,
  };
});
