import { TAG_TYPES } from '@/constants';

export default {
  importedTag: 'Le tag est importé',
  types: {
    [TAG_TYPES.imported]: 'Importé',
    [TAG_TYPES.created]: 'Créé',
  },
  deleteConfirmation: 'Voulez-vous vraiment supprimer ce tag ?\nCette action ne peut pas être annulée. Le tag sera supprimé de toutes les alarmes affectées.'
    + ' | Voulez-vous vraiment supprimer ces tags ?\nCette action ne peut pas être annulée. Les tags seront supprimés de toutes les alarmes concernées.',
};
