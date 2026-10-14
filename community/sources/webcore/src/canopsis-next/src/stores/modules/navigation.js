import { defineStore } from 'pinia';
import { ref } from 'vue';

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
