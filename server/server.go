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

	otpCodeRepository := repository.NewOtpCodeRepository(repository.OCRConfig{
		DB: db.Get(),
	})

	sealabsAccountRepository := repository.NewSealabsAccountRepository(repository.SARConfig{
		DB: db.Get(),
	})

	addressRepository := repository.NewAddressRepository(repository.ARConfig{
		DB: db.Get(),
	})

	merchantVoucherRepository := repository.NewMerchantVoucherRepository(repository.MVRConfig{
		DB: db.Get(),
	})

	productPromo := repository.NewProductPromoRepository(repository.PPRConfig{
		DB: db.Get(),
	})

	checkoutItemRepository := repository.NewCheckoutRepository(repository.CIRConfig{DB: db.Get()})

	courierRepository := repository.NewCourierRepository(repository.CConfig{DB: db.Get()})

	marketplaceVoucherRepository := repository.NewMarketplaceVoucherRepository(repository.MVCRConfig{DB: db.Get()})

	transactionRepository := repository.NewTransactionRepository(repository.TRConfig{DB: db.Get()})

	//SERVICE
	userService := service.NewUserService(service.USConfig{
		UserRepository:    userRepository,
		OtpCodeRepository: otpCodeRepository,
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

	sealabsAccountService := service.NewSealabsAccountService(service.SASConfig{
		SealabsAccountRepository: sealabsAccountRepository,
	})

	merchantVoucherService := service.NewMerchantVoucherService(service.MVSConfig{
		MerchantVoucherRepository: merchantVoucherRepository,
		MerchantRepository:        merchantRepository,
	})

	addressService := service.NewAddressService(service.ASConfig{
		AddressRepository: addressRepository,
	})

	productPromoService := service.NewProductPromoService(service.PPSConfig{
		ProductPromoRepository: productPromo,
		MerchantRepository:     merchantRepository,
		ProductRepository:      productRepository,
	})

	cityService := service.NewCityService(service.CiSConfig{CityRepository: cityRepository})

	checkoutItemService := service.NewCheckoutItemService(
		service.CISConfig{
			CheckoutRepository: checkoutItemRepository,
			ProductRepository:  productRepository,
		})

	courierService := service.NewCourierService(service.CourierSConfig{CourierRepository: courierRepository})

	marketplaceVoucherService := service.NewMarketplaceVoucherService(service.MpVSConfig{MarketplaceVoucherRepository: marketplaceVoucherRepository})

	transactionService := service.NewTransactionService(service.TSConfig{
		TransactionRepository: transactionRepository,
		WalletRepository:      walletRepository,
	})

	r := NewRouter(RouterConfig{
		UserService:               userService,
		ProvinceService:           provinceService,
		MerchantService:           merchantService,
		ProductService:            productService,
		CartService:               cartService,
		MarketplaceService:        marketplaceService,
		WalletService:             walletService,
		CategoryService:           categoryService,
		ProductReviewService:      productReviewService,
		WalletTransactionService:  walletTransactionService,
		WishlistService:           wishlistService,
		SealabsAccountService:     sealabsAccountService,
		MerchantVoucherService:    merchantVoucherService,
		AddressService:            addressService,
		CityService:               cityService,
		CheckoutItemService:       checkoutItemService,
		CourierService:            courierService,
		MarketplaceVoucherService: marketplaceVoucherService,
		ProductPromoService:       productPromoService,
		TransactionService:        transactionService,
	})

	log.Fatal(r.Run())
}
