import { convertSortToRequest } from '@/helpers/entities/shared/query';

export const convertQueryToRequest = ({
  page,
  search,
  itemsPerPage,
  sortBy = [],
  sortDesc = [],
}) => {
  const query = {
    page,
    limit: itemsPerPage,

    ...convertSortToRequest(sortBy, sortDesc),
  };

  if (search) {
    query.search = search;
  }

  return query;
};
