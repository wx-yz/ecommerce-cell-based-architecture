package ecommerce.adservice.controller;

import ecommerce.adservice.model.Ad;
import ecommerce.adservice.model.AdRequest;
import ecommerce.adservice.model.AdResponse;
import ecommerce.adservice.service.AdService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/api/ads")
public class AdController {

    @Autowired
    private AdService adService;

    @PostMapping("/getAds")
    public AdResponse getAds(@RequestBody AdRequest request) {
        List<Ad> ads = adService.getAds(request);
        return new AdResponse(ads);
    }

    @GetMapping("/health")
    public String health() {
        return "Ad Service is healthy";
    }
}