import { TAG_TYPES } from '@/constants';

export default {
  importedTag: 'La balise est importée',
  types: {
    [TAG_TYPES.imported]: 'Importé',
    [TAG_TYPES.created]: 'Créé',
  },
  deleteConfirmation: 'Voulez-vous vraiment supprimer cette balise ?\nCette action ne peut pas être annulée. La balise sera supprimée de toutes les alarmes affectées.'
    + ' | Voulez-vous vraiment supprimer ces balises ?\nCette action ne peut pas être annulée. Les balises seront supprimées de toutes les alarmes concernées.',
};
