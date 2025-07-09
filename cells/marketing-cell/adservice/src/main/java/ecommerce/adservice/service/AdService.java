package ecommerce.adservice.service;

import ecommerce.adservice.model.Ad;
import ecommerce.adservice.model.AdRequest;
import org.springframework.stereotype.Service;

import java.util.ArrayList;
import java.util.List;
import java.util.Random;

@Service
public class AdService {

    private static final Random random = new Random();
    private static final int MAX_ADS_TO_SERVE = 2;

    // Sample ads data
    private static final List<Ad> availableAds = List.of(
        new Ad("ad1", "Shop the latest fashion trends", "https://example.com/ad1", "fashion"),
        new Ad("ad2", "Best deals on electronics", "https://example.com/ad2", "electronics"),
        new Ad("ad3", "Home & Garden essentials", "https://example.com/ad3", "home"),
        new Ad("ad4", "Sports equipment sale", "https://example.com/ad4", "sports"),
        new Ad("ad5", "Books and education", "https://example.com/ad5", "books"),
        new Ad("ad6", "Health and wellness", "https://example.com/ad6", "health"),
        new Ad("ad7", "Travel deals", "https://example.com/ad7", "travel"),
        new Ad("ad8", "Food and beverages", "https://example.com/ad8", "food")
    );

    public List<Ad> getAds(AdRequest request) {
        List<Ad> contextualAds = new ArrayList<>();
        
        // Filter ads based on context if provided
        if (request.getContextKeys() != null && !request.getContextKeys().isEmpty()) {
            for (String context : request.getContextKeys()) {
                availableAds.stream()
                    .filter(ad -> ad.getCategory().toLowerCase().contains(context.toLowerCase()))
                    .limit(MAX_ADS_TO_SERVE)
                    .forEach(contextualAds::add);
            }
        }
        
        // If no contextual ads found, return random ads
        if (contextualAds.isEmpty()) {
            List<Ad> shuffledAds = new ArrayList<>(availableAds);
            java.util.Collections.shuffle(shuffledAds, random);
            return shuffledAds.subList(0, Math.min(MAX_ADS_TO_SERVE, shuffledAds.size()));
        }
        
        return contextualAds;
    }
}