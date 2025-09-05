// Mock product data for the e-commerce storefront

const products = [
  {
    id: 1,
    name: "Premium Wireless Headphones",
    slug: "premium-wireless-headphones",
    description: "High-quality wireless headphones with noise cancellation and premium sound quality. Perfect for music lovers and professionals.",
    price: 299,
    originalPrice: 399,
    category: "Electronics",
    images: [
      "https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=premium%20wireless%20headphones%20black%20modern%20design%20studio%20lighting&image_size=square",
      "https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=wireless%20headphones%20side%20view%20black%20premium%20quality&image_size=square"
    ],
    rating: 4.8,
    reviews: 1247,
    stock: 15,
    variants: {
      colors: ["Black", "White", "Silver"],
      sizes: ["One Size"]
    },
    features: ["Noise Cancellation", "40H Battery", "Quick Charge", "Bluetooth 5.0"],
    material: ["Premium Plastic", "Memory Foam"]
  },
  {
    id: 2,
    name: "Organic Cotton T-Shirt",
    slug: "organic-cotton-t-shirt",
    description: "Comfortable and sustainable organic cotton t-shirt. Made from 100% organic cotton with a relaxed fit.",
    price: 29,
    originalPrice: 39,
    category: "Clothing",
    images: [
      "https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=organic%20cotton%20t-shirt%20white%20minimalist%20design%20flat%20lay&image_size=square",
      "https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=cotton%20t-shirt%20folded%20white%20organic%20fabric&image_size=square"
    ],
    rating: 4.5,
    reviews: 892,
    stock: 25,
    variants: {
      colors: ["White", "Black", "Gray", "Navy"],
      sizes: ["XS", "S", "M", "L", "XL", "XXL"]
    },
    features: ["100% Organic", "Soft Touch", "Breathable", "Machine Washable"],
    material: ["Organic Cotton"]
  },
  {
    id: 3,
    name: "Smart Fitness Watch",
    slug: "smart-fitness-watch",
    description: "Advanced fitness tracking watch with heart rate monitoring, GPS, and smartphone connectivity.",
    price: 199,
    originalPrice: null,
    category: "Electronics",
    images: [
      "https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=smart%20fitness%20watch%20black%20modern%20display%20sports%20design&image_size=square",
      "https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=fitness%20watch%20wrist%20view%20black%20strap%20digital%20display&image_size=square"
    ],
    rating: 4.6,
    reviews: 634,
    stock: 8,
    variants: {
      colors: ["Black", "Silver", "Rose Gold"],
      sizes: ["38mm", "42mm"]
    },
    features: ["Heart Rate Monitor", "GPS Tracking", "Water Resistant", "7-Day Battery"],
    material: ["Aluminum", "Silicone Strap"]
  },
  {
    id: 4,
    name: "Leather Crossbody Bag",
    slug: "leather-crossbody-bag",
    description: "Elegant leather crossbody bag perfect for daily use. Handcrafted with premium genuine leather.",
    price: 89,
    originalPrice: 129,
    category: "Accessories",
    images: [
      "https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=leather%20crossbody%20bag%20brown%20elegant%20design%20studio%20photography&image_size=square",
      "https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=leather%20bag%20brown%20crossbody%20strap%20premium%20quality&image_size=square"
    ],
    rating: 4.7,
    reviews: 456,
    stock: 12,
    variants: {
      colors: ["Brown", "Black", "Tan"],
      sizes: ["Small", "Medium"]
    },
    features: ["Genuine Leather", "Adjustable Strap", "Multiple Pockets", "Handcrafted"],
    material: ["Genuine Leather", "Cotton Lining"]
  },
  {
    id: 5,
    name: "Ceramic Coffee Mug Set",
    slug: "ceramic-coffee-mug-set",
    description: "Beautiful set of 4 ceramic coffee mugs with modern design. Perfect for your morning coffee routine.",
    price: 45,
    originalPrice: null,
    category: "Home & Kitchen",
    images: [
      "https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=ceramic%20coffee%20mug%20set%20white%20modern%20design%20kitchen%20photography&image_size=square",
      "https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=coffee%20mugs%20ceramic%20white%20set%20of%20four%20minimalist&image_size=square"
    ],
    rating: 4.4,
    reviews: 289,
    stock: 20,
    variants: {
      colors: ["White", "Black", "Blue", "Green"],
      sizes: ["12oz", "16oz"]
    },
    features: ["Dishwasher Safe", "Microwave Safe", "Set of 4", "Modern Design"],
    material: ["Ceramic"]
  }
];

const categories = [
  {
    id: 1,
    name: "Electronics",
    slug: "electronics",
    description: "Latest gadgets and electronic devices",
    image: "https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=electronics%20category%20modern%20gadgets%20technology%20devices&image_size=landscape_4_3",
    productCount: 2
  },
  {
    id: 2,
    name: "Clothing",
    slug: "clothing",
    description: "Fashion and apparel for all occasions",
    image: "https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=clothing%20category%20fashion%20apparel%20wardrobe%20style&image_size=landscape_4_3",
    productCount: 1
  },
  {
    id: 3,
    name: "Accessories",
    slug: "accessories",
    description: "Bags, jewelry, and fashion accessories",
    image: "https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=accessories%20category%20bags%20jewelry%20fashion%20items&image_size=landscape_4_3",
    productCount: 1
  },
  {
    id: 4,
    name: "Home & Kitchen",
    slug: "home-kitchen",
    description: "Home decor and kitchen essentials",
    image: "https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=home%20kitchen%20category%20decor%20essentials%20modern%20design&image_size=landscape_4_3",
    productCount: 1
  }
];

// Helper functions
export function getProducts() {
  return products;
}

export function getProduct(slug) {
  return products.find(product => product.slug === slug);
}

export function getProductById(id) {
  return products.find(product => product.id === id);
}

export function getProductsByCategory(categorySlug) {
  return products.filter(product => 
    product.category.toLowerCase().replace(/\s+/g, '-') === categorySlug
  );
}

export function getFeaturedProducts(limit = 4) {
  return products
    .filter(product => product.rating >= 4.5)
    .slice(0, limit);
}

export function getCategories() {
  return categories;
}

export function getCategory(slug) {
  return categories.find(category => category.slug === slug);
}

export function searchProducts(query, filters = {}) {
  let filteredProducts = [...products];
  
  // Text search
  if (query) {
    const searchTerm = query.toLowerCase();
    filteredProducts = filteredProducts.filter(product => 
      product.name.toLowerCase().includes(searchTerm) ||
      product.description.toLowerCase().includes(searchTerm) ||
      product.category.toLowerCase().includes(searchTerm) ||
      (product.features && product.features.some(feature => 
        feature.toLowerCase().includes(searchTerm)
      ))
    );
  }
  
  // Category filter
  if (filters.category && filters.category.length > 0) {
    filteredProducts = filteredProducts.filter(product => 
      filters.category.includes(product.category)
    );
  }
  
  // Price range filter
  if (filters.minPrice !== undefined) {
    filteredProducts = filteredProducts.filter(product => 
      product.price >= filters.minPrice
    );
  }
  
  if (filters.maxPrice !== undefined) {
    filteredProducts = filteredProducts.filter(product => 
      product.price <= filters.maxPrice
    );
  }
  
  // Rating filter
  if (filters.minRating !== undefined) {
    filteredProducts = filteredProducts.filter(product => 
      product.rating >= filters.minRating
    );
  }
  
  // Brand filter (using category as brand for demo)
  if (filters.brand && filters.brand.length > 0) {
    filteredProducts = filteredProducts.filter(product => 
      filters.brand.includes(product.category)
    );
  }
  
  // Features filter
  if (filters.features && filters.features.length > 0) {
    filteredProducts = filteredProducts.filter(product => 
      product.features && filters.features.some(feature => 
        product.features.includes(feature)
      )
    );
  }
  
  // Size filter
  if (filters.size && filters.size.length > 0) {
    filteredProducts = filteredProducts.filter(product => 
      product.variants?.sizes && filters.size.some(size => 
        product.variants.sizes.includes(size)
      )
    );
  }
  
  // Color filter
  if (filters.color && filters.color.length > 0) {
    filteredProducts = filteredProducts.filter(product => 
      product.variants?.colors && filters.color.some(color => 
        product.variants.colors.includes(color)
      )
    );
  }
  
  return filteredProducts;
}

export function sortProducts(products, sortBy) {
  const sortedProducts = [...products];
  
  switch (sortBy) {
    case 'price-low-high':
      return sortedProducts.sort((a, b) => a.price - b.price);
    case 'price-high-low':
      return sortedProducts.sort((a, b) => b.price - a.price);
    case 'rating':
      return sortedProducts.sort((a, b) => b.rating - a.rating);
    case 'reviews':
      return sortedProducts.sort((a, b) => b.reviews - a.reviews);
    case 'newest':
      return sortedProducts.sort((a, b) => b.id - a.id);
    case 'name':
      return sortedProducts.sort((a, b) => a.name.localeCompare(b.name));
    default:
      return sortedProducts;
  }
}

// Get unique values for filters
export function getFilterOptions() {
  const categories = [...new Set(products.map(p => p.category))];
  const brands = categories; // Using categories as brands for demo
  const features = [...new Set(products.flatMap(p => p.features || []))];
  const sizes = [...new Set(products.flatMap(p => p.variants?.sizes || []))];
  const colors = [...new Set(products.flatMap(p => p.variants?.colors || []))];
  const materials = [...new Set(products.flatMap(p => p.material || []))];
  
  const priceRange = {
    min: Math.min(...products.map(p => p.price)),
    max: Math.max(...products.map(p => p.price))
  };
  
  return {
    categories,
    brands,
    features,
    sizes,
    colors,
    materials,
    priceRange
  };
}