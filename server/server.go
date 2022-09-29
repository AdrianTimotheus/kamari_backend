package server

import (
	"log"

	"kamari.com/db"
	"kamari.com/repository"
	"kamari.com/service"
)

func Init() {
	//REPO
	userRepository := repository.NewUserRepository(repository.URConfig{
		DB: db.Get(),
	})

	merchantRepository := repository.NewMerchantRepository(repository.MRConfig{
		DB: db.Get(),
	})

	productRepository := repository.NewProductRepository(repository.ProdConfig{
		DB: db.Get(),
	})
	productVariantRepository := repository.NewProductVariantRepository(repository.PVConfig{
		DB: db.Get(),
	})
	categoryRepositoy := repository.NewCategoryRepository(repository.CaRConfig{
		DB: db.Get(),
	})
	orderDetailRepository := repository.NewOrderDetailRepository(repository.ODRConfig{
		DB: db.Get(),
	})
	merchantFollowerRepository := repository.NewMFRepository(repository.MFConfig{
		DB: db.Get(),
	})
	productReviewRepository := repository.NewProductReviewRepository(repository.PRRConfig{
		DB: db.Get(),
	})

	cityRepository := repository.NewCityRepository(repository.CRConfig{
		DB: db.Get(),
	})

	provinceRepository := repository.NewProvinceRepository(repository.PRConfig{
		DB: db.Get(),
	})

	marketplaceRepository := repository.NewMarketplaceRepository(repository.MpRConfig{
		DB: db.Get(),
	})

	cartRepository := repository.NewCartRepository(repository.CRConfig{
		DB: db.Get(),
	})

	walletRepository := repository.NewWalletRepository(repository.WConfig{
		DB: db.Get(),
	})

	walletTransactionRepository := repository.NewWalletTransaction(repository.WTConfig{
		DB: db.Get(),
	})

	wishlistRepository := repository.NewWishlistRepository(repository.WLConfig{
		DB: db.Get(),
	})

	//SERVICE
	userService := service.NewUserService(service.USConfig{
		UserRepository: userRepository,
	})
	productService := service.NewProductService(service.PConfig{
		ProductRepository:        productRepository,
		ProductVariantRepository: productVariantRepository,
		OrderDetailsRepository:   orderDetailRepository,
		CityRepository:           cityRepository,
		MerchantRepository:       merchantRepository,
	})
	categoryService := service.NewCategoryService(service.CaSConfig{
		CategoryRepository: categoryRepositoy,
	})
	merchantService := service.NewMerchantService(service.MSConfig{
		MerchantRepository:         merchantRepository,
		CityRepository:             cityRepository,
		ProductRepository:          productRepository,
		MerchantFollowerRepository: merchantFollowerRepository,
	})
	productReviewService := service.NewProductReviewService(service.PRSConfig{
		ProductReviewRepository:  productReviewRepository,
		OrderDetailRepository:    orderDetailRepository,
		ProductVariantRepository: productVariantRepository,
	})

	provinceService := service.NewProvinceService(service.PSConfig{
		ProvinceRepository: provinceRepository,
	})

	marketplaceService := service.NewMarketplaceService(service.MpSConfig{
		MarketplaceRepository: marketplaceRepository,
	})

	cartService := service.NewCartService(service.CSConfig{
		CartRepository: cartRepository,
	})

	walletService := service.NewWalletService(service.WSConfig{
		WalletRepository:            walletRepository,
		WalletTransactionRepository: walletTransactionRepository,
	})

	walletTransactionService := service.NewWalletTransactionService(service.WTSConfig{
		WalletTransactionRepository: walletTransactionRepository,
		WalletRepository:            walletRepository,
	})

	wishlistService := service.NewWishlistService(service.WLConfig{
		WishlistRepository: wishlistRepository,
	})

	r := NewRouter(RouterConfig{
		UserService:              userService,
		ProvinceService:          provinceService,
		MerchantService:          merchantService,
		ProductService:           productService,
		CartService:              cartService,
		MarketplaceService:       marketplaceService,
		WalletService:            walletService,
		CategoryService:          categoryService,
		ProductReviewService:     productReviewService,
		WalletTransactionService: walletTransactionService,
		WishlistService:          wishlistService,
	})

	log.Fatal(r.Run())
}
