export type PriceBreakdown = {
  input_per_million: number;
  output_per_million: number;
  cache_read_per_million: number;
  cache_write_per_million: number;
  unit_price: number;
  multiplier: number;
};

export type CatalogModelItem = {
  model: string;
  display_name: string;
  provider: string;
  endpoint_type: string;
  billing_type: number;
  billing_label: string;
  tags: string[];
  prices: PriceBreakdown;
  prices_applied: PriceBreakdown;
};

export type TokenGroupItem = {
  slug: string;
  name: string;
  multiplier: number;
};

export type CatalogListResponse = {
  success: boolean;
  data: {
    items: CatalogModelItem[];
    total: number;
    page: number;
    page_size: number;
    multiplier: number;
    token_group: string;
  };
  meta: {
    providers: string[];
    token_groups: TokenGroupItem[];
  };
};
