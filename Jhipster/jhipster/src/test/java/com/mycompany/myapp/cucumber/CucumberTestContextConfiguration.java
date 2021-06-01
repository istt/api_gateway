package com.mycompany.myapp.cucumber;

import com.mycompany.myapp.JhipsterApp;
import io.cucumber.spring.CucumberContextConfiguration;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.web.WebAppConfiguration;

@CucumberContextConfiguration
@SpringBootTest(classes = JhipsterApp.class)
@WebAppConfiguration
public class CucumberTestContextConfiguration {}
